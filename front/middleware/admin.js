export default defineNuxtRouteMiddleware(async (to, from) => {
    const { onLogout, getToken } = useApollo()
    let token = await getToken()


    let decodedToken = atob(token.split('.')[1])

    let tokenParsed = JSON.parse(decodedToken)

    if (tokenParsed.role !== 9) {
        return navigateTo('/')
    }
})