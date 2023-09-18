<template>
    <div class="hero min-h-screen bg-base-200" style="background-image: url(https://daisyui.com/images/stock/photo-1507358522600-9f71e620c44e.jpg);">
        <div class="hero-overlay bg-opacity-60"></div>
        <div class="hero-content flex-col lg:flex-row-reverse text-neutral-content">
            <div class="text-center lg:text-left">
                <h1 class="text-5xl font-bold">Mot de passe oublié</h1>
            </div>
            <div v-if="success" class="card flex-shrink-0 w-full max-w-sm shadow-2xl bg-success">
                <div class="card-body text-black text-center">
                    <p>Un email contenant les étapes à suivre a été envoyé à l'adresse indiquée. </p>
                </div>
            </div>
            <div v-else class="card flex-shrink-0 w-full max-w-sm shadow-2xl bg-base-100">
                <div class="card-body">
                    <div class="form-control">
                        <label class="label">
                            <span class="label-text">Renseigner l'email dont vous souhaitez changer le mot de passe</span>
                        </label>
                        <input v-model="email" type="text" placeholder="email@gmail.com" class="input input-bordered text-black" />
                    </div>
                    <div class="form-control mt-6">
                        <button @click="onRegister" class="btn btn-primary">
                            <span class="loading loading-dots loading-sm" v-if="actionInProgress"></span>
                            <span v-else>Envoyer la demande</span>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
definePageMeta({
  layout: "custom",
});
</script>

<script>
export default {
    name: "AskResetPasswordPage",
    data() {
        return {
            email: '',
            actionInProgress: false,
            success: false
        }
    },
    methods: {
        async onRegister () {
            if (this.actionInProgress) return

            this.actionInProgress = true

            let query = gql`
                mutation {
                    askResetPassword(input: {email: "${this.email}"}) {
                        email
                    }
                }
            `

            try {
                await useAsyncQuery(query)
                this.success = true
            } catch(e) {
            }

            this.actionInProgress = false
        }
    }
}
</script>