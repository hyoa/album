export default defineNuxtRouteMiddleware(async (to, from) => {
    const { onLogout, getToken } = useApollo()
    let token = await getToken()
    
    if (token === null && token !== "") {
        console.log(token)
        console.log("redirect")
        return navigateTo('/login')
    }

    let decodedToken = atob(token.split('.')[1])

    let tokenParsed = JSON.parse(decodedToken)


    const now = Date.now()

    if (now > tokenParsed.exp * 1000) { 
        localStorage.removeItem('albumToken')
        localStorage.removeItem('albumTokenExp')
        token = null

        onLogout()
        return navigateTo('/login')
    }
})