<template>
  <AdminLayout>
    <template v-slot:title>
      <PageTitle title="Créer un album" icon="regular/plus-square" color="bg-indigo-500"/>
    </template>
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
import { graphql } from '../../../utils/axiosHelper'

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

      const query = `
        mutation {
          createAlbum(input: {title: "${this.title}", author: "${this.$store.state.token.name}", description: "${this.description}", private: ${!!this.isPrivate}}) {
            title
          }
        }
      `

      graphql(query, 'v3')
        .then(() => {
          this.$notify({ group: 'success', text: 'L\'album a été créé avec succès' })
        })
        .catch(message => {
          this.$notify({ group: 'error', text: message })
        })
        .finally(() => {
          this.formStatus.add = 'ready'
        })
    }
  }
}
</script>
