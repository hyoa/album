<template>
    <div class="px-4">
        <h1 class="text-xl">{{ title }}</h1>
    </div>
</template>

<script>
    export default {
        name: "PageTitle",
        props: {
            title: {
                type: String,
                required: true
            }
        }
    }
</script>