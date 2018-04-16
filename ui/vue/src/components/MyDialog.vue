<template>

  <!-- modal dialog -->
  <b-modal id="modal1" ref="modal1" title="Input Transaksi" @shown="reset" @ok="handleOk">

    <!-- input deskripsi -->
    <b-form-group label="Deskripsi">
      <b-form-input type="text" v-model="obj.deskripsi" required></b-form-input>
    </b-form-group>

    <!-- input nilai -->
    <b-form-group label="Nilai">
      <!-- <b-form-input type="text" v-model.number="obj.nilai" required></b-form-input> -->
      <money v-model="obj.nilai" class="form-control text-right"></money>
    </b-form-group>
  </b-modal>
</template>

<script>
export default {
  props: ['input'],
  data () {
    return {
      obj: {
        nilai: 0,
        deskripsi: ''
      }
    }
  },
  methods: {
    handleOk (evt) {
      evt.preventDefault()
      if (this.input == null) {
        this.$http.post('/transaksi', this.obj).then((response) => {
          this.reset()
          this.$refs.modal1.hide()
          this.$emit('closed', 'created')
        }).catch((error) => {
          this.$swal('Error', error.response.data.message, 'error')
        })
      } else {
        this.$http.put('/transaksi/' + this.input.id, this.obj).then((response) => {
          this.reset()
          this.$refs.modal1.hide()
          this.$emit('closed', 'updated')
        }).catch((error) => {
          this.$swal('Error', error.response.data.message, 'error')
        })
      }
    },
    reset () {
      if (this.input == null) {
        this.obj = {
          nilai: 0,
          deskripsi: ''
        }
      } else {
        this.obj = this.input
      }
    }
  }
}
</script>

<style>
.text-right {
  text-align: right
}
</style>
