<template>
    <div class="hero min-h-screen bg-base-200" style="background-image: url(https://daisyui.com/images/stock/photo-1507358522600-9f71e620c44e.jpg);">
        <div class="hero-overlay bg-opacity-60"></div>
        <div class="hero-content flex-col lg:flex-row-reverse text-neutral-content">
            <div class="text-center lg:text-left">
                <h1 class="text-5xl font-bold">Changer de mot de passe</h1>
            </div>
            <div v-if="success" class="card flex-shrink-0 w-full max-w-sm shadow-2xl bg-success">
                <div class="card-body text-black text-center">
                    <p>Votre mot de passe a été modifié</p>
                    <NuxtLink class="text-blue-500 underline" to="/login" >Se connecter</NuxtLink>
                </div>
            </div>
            <div v-else class="card flex-shrink-0 w-full max-w-sm shadow-2xl bg-base-100">
                <div class="card-body">
                    <div class="text-red-500" v-if="errorMessage !== ''">
                        {{ errorMessage }}
                    </div>
                    <div class="form-control">
                        <label class="label">
                            <span class="label-text">Mot de passe</span>
                        </label>
                        <input v-model="password" type="password" placeholder="*******" class="input input-bordered text-black" />
                    </div>
                    <div class="form-control">
                        <label class="label">
                            <span class="label-text">Confirmation du mot de passe</span>
                        </label>
                        <input v-model="checkPassword" type="password" placeholder="*******" class="input input-bordered text-black" />
                    </div>
                    <div class="form-control mt-6">
                        <button @click="onReset" class="btn btn-primary">
                            <span class="loading loading-dots loading-sm" v-if="actionInProgress"></span>
                            <span v-else>Mettre à jour</span>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
definePageMeta({
  layout: "none",
});
</script>

<script>
export default {
    name: "ResetPasswordPage",
    data() {
        return {
            password: '',
            checkPassword: '',
            actionInProgress: false,
            success: false,
            token: null,
            errorMessage: ''
        }
    },
    created () {
        this.token = this.$route.query.token
    },
    methods: {
        async onReset () {
            this.errorMessage = ''
            if (this.actionInProgress) return

            this.actionInProgress = true

            let query = gql`
                mutation {
                    resetPassword(input: {password: "${this.password}", passwordCheck: "${this.checkPassword}", tokenValidation: "${this.token}"}) {
                        email
                    }
                }
            `

            try {
                const { data: { _rawValue: { resetPassword } } } = await useAsyncQuery(query)
                this.success = true
            } catch(e) {
                this.errorMessage = 'Une erreur est survenue.'
            }

            this.actionInProgress = false
        }
    }
}
</script>