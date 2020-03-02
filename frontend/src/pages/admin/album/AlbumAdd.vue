<template>
  <AdminLayout>
    <PageTitle title="Créer un album" />
    <form @submit.prevent="onSubmit">
      <InputSimple v-model="title" id="title" placeholder="Mon super album" label="Titre" type="text" />
      <TextareaSimple v-model="description" id="description" placeholder="Raconte moi une histoire..." label="Description" />
      <CheckboxSimple v-model="isPrivate" label="Privé" />
      <SimpleAnimateButton :status="formStatus.add">
        Enregistrer
      </SimpleAnimateButton>
    </form>
  </AdminLayout>
</template>

<script>
import { post } from '../../../utils/axiosHelper'
import errorHelper from '../../../utils/errorHelper'

import AdminLayout from '../../../components/layout/AdminLayout'
import InputSimple from '../../../components/form/default/InputSimple'
import TextareaSimple from '../../../components/form/default/TextareaSimple'
import CheckboxSimple from '../../../components/form/default/CheckboxSimple'
import PageTitle from '../../../components/admin/PageTitle'
import SimpleAnimateButton from '../../../components/form/button/SimpleAnimateButton'

export default {
  name: 'AdminAlbumAdd',
  components: { SimpleAnimateButton, PageTitle, CheckboxSimple, TextareaSimple, InputSimple, AdminLayout },
  data () {
    return {
      title: '',
      description: '',
      isPrivate: false,
      formStatus: {
        add: 'ready'
      }
    }
  },
  methods: {
    onSubmit () {
      if (!this.title) {
        return
      }
      this.formStatus.add = 'pending'

      const data = {
        title: this.title,
        description: this.description,
        private: !!this.isPrivate
      }

      post(`album`, data)
        .then(() => {
          this.$notify({ group: 'success', text: 'L\'album a été créé avec succès' })
        })
        .catch(({ response }) => {
          let code = null
          try {
            code = response.data.code
          } catch (e) {
            code = 999
          }

          this.$notify({ group: 'error', text: errorHelper(code) })
        })
        .finally(() => {
          this.formStatus.add = 'ready'
        })
    }
  }
}
</script>
