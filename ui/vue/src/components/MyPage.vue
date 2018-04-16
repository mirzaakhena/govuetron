<template>
  <div>

    <!-- button utk buka dialog input -->
    <b-btn variant="success" v-b-modal.modal1>Transaksi Baru</b-btn>

    <!-- dialog input -->
    <my-dialog :input="selectedItem" @closed="dialogClosed"></my-dialog>

    <!-- tabel -->
    <b-table class="gap-top" small striped hover :items="items" :fields="fields">

      <!-- field tanggal -->
      <template slot="tanggal" slot-scope="data" >
        {{data.item.tanggal | dateformat}}
      </template>

      <!-- field nilai -->
      <template slot="nilai" slot-scope="data" >
        {{data.item.nilai | digitgroup}}
      </template>

      <!-- field delete -->
      <template slot="tindakan" slot-scope="data">
        <b-button-group>
          <b-button v-b-modal.modal1 @click="selectedItem=data.item" size="sm" variant="primary"><icon name="pencil"></icon> Ubah</b-button>
          <b-button @click="deleteItem(data.item)" size="sm" variant="danger"><icon name="remove"></icon> Hapus</b-button>
        </b-button-group>
      </template>

    </b-table>

  </div>

</template>

<script>

import MyDialog from './MyDialog'

export default {
  components: {
    MyDialog
  },
  created () {
    this.reload()
    this.$options.sockets.onmessage = data => this.reload()
  },
  data () {
    return {
      selectedItem: null,
      items: [],
      fields: ['tanggal', 'deskripsi', 'nilai', 'tindakan']
    }
  },
  methods: {
    reload () {
      this.$http.get('/transaksi').then(response => {
        this.items = response.data.data
      }).catch(error => {
        this.$swal('Error', error.response.data.message, 'error')
      })
    },
    dialogClosed (message) {
      this.selectedItem = null
      if (message === 'created') {
        this.reload()
        this.$socket.sendObj({changedDataPage: 'item created'})
      } else if (message === 'updated') {
        this.$socket.sendObj({changedDataPage: 'item updated'})
      }
      // delete this.$options.sockets.onmessage
    },
    deleteItem (item) {
      this.$swal({
        title: 'Are you sure?',
        text: 'You will delete the item',
        type: 'warning',
        showCancelButton: true
      }).then(result => {
        if (result.value) {
          this.$http.delete('/transaksi/' + item.id).then(response => {
            this.reload()
            this.$socket.sendObj({changedDataPage: 'item deleted'})
          }).catch(error => {
            this.$swal('Error', error.response.data.message, 'error')
          })
        }
      })
    }
  }
}
</script>

<style>

.gap-top {
  margin-top: 20px
}

table td {
  vertical-align: middle !important;
}

.btn .fa-icon {
  vertical-align: sub;
}

</style>
