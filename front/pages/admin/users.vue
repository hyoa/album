<template>
    <section class="px-4">
        <div class="text-sm breadcrumbs">
            <ul>
                <li><NuxtLink to="/admin">Admin</NuxtLink></li> 
                <li class="font-bold">Utilisateurs</li>
            </ul>
        </div>
    </section>
    <section class="px-4 mt-4">
        <h2 class="text-gray-600">Inviter un utilisateur</h2>
        <div class="w-full mt-2">
            <input v-model="emailToInvite" type="text" placeholder="email@email.com" class="input input-bordered w-full max-w-xs" />
        </div>
        <div class="w-full text-center mt-2">
            <button @click="onInvite" class="btn btn-block btn-primary">Inviter</button>
        </div>
    </section>
    <div class="divider px-4"></div>
    <section class="px-4">
        <h2 class="text-gray-600">Utilisateurs enregistrés</h2>
        <div class="overflow-x-auto">
            <table class="table">
                <!-- head -->
                <thead>
                <tr>
                    <th>Name</th>
                    <th>Role</th>
                    <th>Actions</th>
                </tr>
                </thead>
                <tbody>
                    <tr v-for="(user, index) in users" :key="index">
                        <td>
                            <p>{{ user.name }}</p>
                            <p>{{ user.email }}</p>
                        </td>
                        <td>{{ user.role }}</td>
                        <td v-if="user.role === 'UNIDENTIFIED'">
                            <button @click="onActivate(user.email)" class="btn btn-sm btn-secondary">Activer</button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </section>
    <div class="toast toast-end" v-if="alert.type !== null">
        <div class="alert" :class="alert.type == 'success' ? 'alert-success': 'alert-error'">
            <span>{{ alert.message }}</span>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'UsersList',
        data() {
            return {
                users: [],
                alert: {
                    type: null,
                    title: null,
                    message: null
                },
                emailToInvite: '',
                formStatus: {
                    invite: 'ready'
                }
            }
        },
        async created() {
            const query = gql`
                query {
                    users: users {
                        name
                        email
                        role
                    }
                }
            `

            const { data: { _rawValue: { users } } } = await useAsyncQuery(query, {})
            this.users = users
        },
        methods: {
            async onActivate(email) {
                this.alert.type = null

                const mutation = gql`
                    mutation updateUser($input: UpdateInput!) {
                        user: updateUser(input: $input) {
                            role
                        }
                    }
                `

                const variables = {
                    input: {
                        email,
                        role: 'NORMAL'
                    }
                }

                const { data: { _rawValue: { user } } } = await useAsyncQuery(mutation, variables)
                
                if (user.role === 'NORMAL') {
                    this.alert.type = 'success'
                    this.alert.title = 'Utilisateur activé'
                    this.alert.message = 'L\'utilisateur a bien été activé.'
                } else {
                    this.alert.type = 'error'
                    this.alert.title = 'Erreur'
                    this.alert.message = 'Une erreur est survenue lors de l\'activation de l\'utilisateur.'
                }
            },
            async onInvite() {
                this.alert.type = null

                if (this.emailToInvite === '') {
                    return
                }

                const mutation = gql`
                    mutation invite($input: InviteInput!) {
                        invite(input: $input) {
                            email
                        }
                    }
                `

                const variables = {
                    input: {
                        email: this.emailToInvite
                    }
                }

                const { data: { _rawValue: { invite } } } = await useAsyncQuery(mutation, variables)
                
                if (invite.email === this.emailToInvite) {
                    this.alert.type = 'success'
                    this.alert.title = 'Invitation envoyée'
                    this.alert.message = 'L\'invitation a bien été envoyée.'
                } else {
                    this.alert.type = 'error'
                    this.alert.title = 'Erreur'
                    this.alert.message = 'Une erreur est survenue lors de l\'envoi de l\'invitation.'
                }
            }
        }
    }
</script>
