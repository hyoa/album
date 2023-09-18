<template>
    <div>
        <div class="navbar bg-base-100">
            <div class="flex-1">
                <NuxtLink to="/" class="btn btn-ghost normal-case text-xl">Pauline & Jules</NuxtLink>
            </div>
            <div class="flex-none">

                <div class="dropdown dropdown-end">
                    <label tabindex="0" class="btn btn-ghost btn-circle avatar">
                        <div class="w-10 rounded-full">
                        <img :src="`https://ui-avatars.com/api/?rounded=true&name=${username}`" />
                        </div>
                    </label>
                    <ul tabindex="0" class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52">
                        <li v-if="role === 9">
                            <NuxtLink to="/admin">
                                Administation
                            </NuxtLink>
                        </li>
                        <li @click="logout">
                            <span>
                                Se déconnecter
                            </span>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
        <div class="mt-2">
            <slot />
        </div>
    </div>
</template>

<script>
const { onLogout, getToken } = useApollo()

export default {
    name: "DefaultLayout",
    components: {
    },
    data() {
        return {
            username: '',
            role: null
        }
    },
    async created() {
        let token = await getToken()
        let decodedToken = atob(token.split('.')[1])
        let tokenParsed = JSON.parse(decodedToken)

        this.username = tokenParsed.name
        this.role = tokenParsed.role
    },
    methods: {
        async logout() {
            await onLogout()
            return navigateTo('/login')
        }
    }
}
</script>