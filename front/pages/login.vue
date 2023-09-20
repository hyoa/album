<template>
    <div class="hero min-h-screen bg-base-200" style="background-image: url(https://daisyui.com/images/stock/photo-1507358522600-9f71e620c44e.jpg);">
        <div class="hero-overlay bg-opacity-60"></div>
        <div class="hero-content flex-col lg:flex-row-reverse text-neutral-content">
            <div class="text-center lg:text-left">
                <h1 class="text-5xl font-bold">Se connecter</h1>
            </div>
            <div class="card flex-shrink-0 w-full max-w-sm shadow-2xl bg-base-100">
                <div class="card-body">
                    <div class="text-red-500" v-if="errorMessage !== ''">
                        {{ errorMessage }}
                    </div>
                    <div class="form-control">
                        <label class="label">
                            <span class="label-text">Email</span>
                        </label>
                        <input v-model="email" type="text" placeholder="email@gmail.com" class="input input-bordered text-black" />
                    </div>
                    <div class="form-control">
                        <label class="label">
                            <span class="label-text">Mot de passe</span>
                        </label>
                        <input v-model="password" type="password" placeholder="*******" class="input input-bordered text-black" />
                        <label class="label">
                            <NuxtLink to="/askResetPassword" class="label-text-alt link link-hover">Mot de passe oublié ?</NuxtLink>
                        </label>
                    </div>
                    <div class="form-control mt-6">
                        <button @click="login" class="btn btn-primary">
                            <span class="loading loading-dots loading-sm" v-if="actionInProgress"></span>
                            <span v-else>Se connecter</span>
                        </button>
                    </div>
                    <div class="text-center">
                        <NuxtLink class="text-blue-500 underline" to="/register" >S'inscrire</NuxtLink>
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
import graphql from '../internal/graphql';
const { onLogin } = useApollo()

export default {
    name: 'Login',
    data() {
        return {
            email: '',
            password: '',
            errorMessage: '',
            actionInProgress: false,
        }
    },
    methods: {
        async login() {
            if (this.actionInProgress) return;

            this.actionInProgress = true;

            this.errorMessage = '';
            let query = gql`
            query {
                auth: auth(input: {email: "${this.email}", password: "${this.password}"}) {
                token
                }
            }
            `

            try {
                const { data: { _rawValue: { auth } } } = await useAsyncQuery(query)
                await onLogin(auth.token)

                navigateTo('/')
            } catch (e) {
                this.errorMessage = e;
            }

            this.actionInProgress = false;
        },
    },
}
</script>