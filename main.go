package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/olahol/melody"
	"github.com/spf13/viper"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/satori/go.uuid"
)

//--------

type Application struct {
	DB *gorm.DB
}

type Transaksi struct {
	ID        string    `json:"id"`
	Deskripsi string    `json:"deskripsi"`
	Tanggal   time.Time `json:"tanggal"`
	Nilai     float64   `json:"nilai"`
}

type StandartResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//--------

func (app Application) simpanData(c *gin.Context) {

	// siapkan requestnya
	var request struct {
		Deskripsi string  `json:"deskripsi"`
		Nilai     float64 `json:"nilai"`
	}

	// bind objectnya
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, StandartResponse{"Gagal Binding Object", nil})
		return
	}

	// validasi input
	if request.Nilai == 0 {
		c.JSON(http.StatusBadRequest, StandartResponse{"Nilai tidak boleh 0", nil})
		return
	}

	tx := app.DB.Begin()

	// siapkan objeknya
	obj := Transaksi{
		ID:        uuid.NewV4().String(),
		Deskripsi: request.Deskripsi,
		Tanggal:   time.Now(),
		Nilai:     request.Nilai,
	}

	// simpan objeknya
	if err := tx.Create(&obj).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, StandartResponse{"Gagal Insert", obj})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, StandartResponse{"Sukses Insert", obj})
}

func (app Application) ambilSemuaData(c *gin.Context) {

	// ambil semua transaksi
	var listTrx []Transaksi
	app.DB.Find(&listTrx)

	c.JSON(http.StatusOK, StandartResponse{"Sukses Ambil Data", listTrx})
}

func (app Application) ubahData(c *gin.Context) {

	// ambil transaksi_id nya
	transaksiID := c.Param("transaksiID")

	// siapkan requestnya
	var request struct {
		Deskripsi string  `json:"deskripsi"`
		Nilai     float64 `json:"nilai"`
	}

	// bind objectnya
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, StandartResponse{"Gagal Binding Object", nil})
		return
	}

	// validasi input
	if request.Nilai == 0 {
		c.JSON(http.StatusBadRequest, StandartResponse{"Nilai tidak boleh 0", nil})
		return
	}

	// ambil transaksi by id
	var trx Transaksi
	app.DB.Where("ID = ?", transaksiID).First(&trx)

	// ubah nilai2 nya
	trx.Nilai = request.Nilai
	trx.Deskripsi = request.Deskripsi

	// simpan kembali
	app.DB.Save(&trx)

	c.JSON(http.StatusOK, StandartResponse{"Sukses Ubah Data", trx})
}

func (app Application) hapusData(c *gin.Context) {

	// ambil transaksi_id nya
	transaksiID := c.Param("transaksiID")

	// ambil transaksi by id
	app.DB.Delete(Transaksi{}, "ID = ?", transaksiID)

	c.JSON(http.StatusOK, StandartResponse{"Sukses Hapus Data", nil})
}

//--------

//SocketEvent event
type SocketEvent struct {
	Event string
	Data  interface{}
}

//Backend Configuration returned by loadConfig
type configuration struct {
	name   string //Application Name
	server string
	port   string
}

//Globals
var done = make(chan bool, 1)      //Wait for Shutdown signal over websocket
var upgrader = websocket.Upgrader{ //Upgrader for websockets
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Subprotocols:    []string{"p0", "p1"},
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Loads configuration from file
// or inits values with default values
func loadConfig() configuration {
	viper.SetConfigName("config")

	// Paths to search for a config file
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("No configuration file loaded - using defaults")
		color.Unset()
	}

	// default values
	viper.SetDefault("name", "")
	viper.SetDefault("server", "localhost")
	viper.SetDefault("port", "9109")

	// Write all params to stdout
	color.Set(color.FgGreen)
	fmt.Println("Loaded Configuration:")
	color.Unset()

	// Print config
	keys := viper.AllKeys()
	for i := range keys {
		key := keys[i]
		fmt.Println(key + ":" + viper.GetString(key))
	}
	fmt.Println("---")

	return configuration{
		name:   viper.GetString("name"),
		server: viper.GetString("server"),
		port:   viper.GetString("port")}
}

//Handles msgs to communicate with nodejs electron for rampup & shutdown
func socket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		var event SocketEvent
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("ElectronSocket: [err]", err)
			break
		}

		//Handle Message
		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Println("Unmashal: ", err)
			break
		}
		log.Printf("ElectronSocket: [received] %+v", event)

		//Shutdown Event
		if event.Event == "shutdown" {
			switch t := event.Data.(type) {
			case bool:
				if t {
					done <- true
				}
			}
		}
	}
}

func main() {
	config := loadConfig()

	var addr = config.server + ":" + config.port
	http.HandleFunc("/ui", socket)    //Endpoint for Electron startup/teardown
	go http.ListenAndServe(addr, nil) //Start websockets in goroutine

	log.Printf("Starting Electron...")
	cmd := exec.Command("electron", "app")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	color.Set(color.FgGreen)
	log.Printf("%s succesfully started", config.name)
	color.Unset()

	//----------

	r := gin.Default()
	r.Use(cors.Default())

	m := melody.New()

	// buka koneksi
	db, err := gorm.Open("sqlite3", "transaksi.db")
	if err != nil {
		panic("failed to connect database")
	}

	// nanti ditutup
	defer db.Close()

	// bikin tabel
	db.AutoMigrate(&Transaksi{})

	app := Application{
		DB: db,
	}

	r.POST("/transaksi", app.simpanData)
	r.GET("/transaksi", app.ambilSemuaData)
	r.PUT("/transaksi/:transaksiID", app.ubahData)
	r.DELETE("/transaksi/:transaksiID", app.hapusData)

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.BroadcastFilter(msg, func(q *melody.Session) bool {
			return q.Request.URL.Path == s.Request.URL.Path
		})
	})

	go r.Run(":8081")
	//---------

	<-done //Wait for shutdown signal
	color.Set(color.FgGreen)
	log.Printf("Shutting down...")
	color.Unset()
}
