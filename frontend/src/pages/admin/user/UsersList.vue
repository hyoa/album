<template>
  <AdminLayout>
    <template v-slot:title>
       <PageTitle :title="$t('admin.userList.title')" icon="list" color="bg-pink-500"/>
    </template>
    <Alert v-if="alert.type" :type="alert.type" :message="alert.message" :title="alert.title" />
    <section>
      <h2>{{ $t('admin.userList.sectionInvitation.title') }}</h2>
      <InputForm v-model="emailsToInvite" placeholder="email1@gmail.com,email2@gmail.com" id="emails" />
      <ButtonForm @click="onInvite" :status="formStatus.invite">{{ $t('admin.userList.sectionInvitation.submit') }}</ButtonForm>
    </section>
    <section>
      <table class="w-full mt-6">
        <tr class="text-left bg-light-primary">
          <th class="p-1">{{ $t('admin.userList.sectionList.table.name') }}</th>
          <th class="p-1">{{ $t('admin.userList.sectionList.table.role') }}</th>
          <th class="p-1">{{ $t('admin.userList.sectionList.table.action') }}</th>
        </tr>
        <tr class="hover:bg-lighter-primary bg-white" v-for="user of users" :key="user.email">
          <td class="py-2 px-1">
            <div>{{ user.name }}</div>
            <div class="text-xs">{{ user.email | truncate(20) }}</div>
          </td>
          <td class="py-2 px-1">{{ user.role }}</td>
          <td class="py-2 px-1">
            <button
              class="bg-transparent hover:bg-blue text-blue-dark font-semibold hover:text-white py-1 px-2 border border-blue hover:border-transparent rounded"
              @click.stop="onActivate(user.email)"
              v-if="user.role === 0"
            >
              {{ $t('admin.userList.sectionList.table.activate') }}
            </button>
          </td>
        </tr>
      </table>
    </section>

  </AdminLayout>
</template>

<script>
import { get, post } from '../../../utils/axiosHelper'
import errorHelper from '../../../utils/errorHelper'
import AdminLayout from '../../../components/layout/AdminLayout'
import InputForm from '../../../components/form/default/InputSimple'
import ButtonForm from '../../../components/form/button/SimpleAnimateButton'
import Alert from '../../../components/alerts/Alert'
import PageTitle from '../../../components/admin/PageTitle'
export default {
  name: 'UsersList.vue',
  components: { PageTitle, Alert, AdminLayout, InputForm, ButtonForm },
  data () {
    return {
      users: [],
      alert: {
        type: null,
        title: null,
        message: null
      },
      emailsToInvite: '',
      formStatus: {
        invite: 'ready'
      }
    }
  },
  async created () {
    const res = await get('users')
    this.users = res.data
  },
  methods: {
    onActivate (email) {
      post('user/activate', { email, role: 1 })
        .then(() => {
          this.$notify({ group: 'success', text: this.$t('admin.userList.notify.accountActivate') })
        })
        .catch(({ response }) => {
          let code = null
          try {
            code = response.data.code
          } catch (e) {
            code = 999
          }

          this.$notify({ group: 'error', text: this.$t(errorHelper(code)) })
        })
    },
    onInvite () {
      this.formStatus.invite = 'pending'

      if (this.emailsToInvite === '') {
        return
      }

      post('users/invite', { emails: this.emailsToInvite })
        .then(() => {
          this.$notify({ group: 'success', text: this.$t('admin.userList.notify.invitationSend') })
        })
        .catch(({ response }) => {
          let code = null
          try {
            code = response.data.code
          } catch (e) {
            code = 999
          }

          this.$notify({ group: 'error', text: this.$t(errorHelper(code)) })
        })
        .finally(() => {
          this.formStatus.invite = 'ready'
        })
    }
  }
}
</script>
