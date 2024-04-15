
<script>
    import getStream from '../views/StreamView.vue'
    export default {
        props: ['page_number', 'number_of_pages'],
    }
</script>
<template>
    <div class="btn-group btn-group-toggle" data-toggle="buttons">
        <label v-if="page_number!=0" class="btn btn-dark" @click="getStream(page_number-1)">
            previous page
        </label>
        <label v-if="page_number!=0 && page_number-3>0" class="btn btn-dark" @click="getStream(page_number-3)">
            {{ page_number-3 }}
        </label>
        <label v-if="page_number!=0 && page_number-2>0" class="btn btn-dark" @click="getStream(page_number-2)">
            {{ page_number-2 }}
        </label>
        <label v-if="page_number!=0 && page_number-1>0" class="btn btn-dark" @click="getStream(page_number-1)">
            {{ page_number-1 }}
        </label>
        <label class="btn btn-dark active">{{ page_number }} </label>
        <label v-if="page_number+1<=number_of_pages" class="btn btn-dark" @click="getStream(page_number+1)">
            {{ page_number+1 }}
        </label>
        <label v-if="page_number+2<=number_of_pages" class="btn btn-dark" @click="getStream(page_number+2)">
            {{ page_number+2 }}
        </label>
        <label v-if="page_number+3<=number_of_pages" class="btn btn-dark" @click="getStream(page_number+3)">
            {{ page_number+3 }}
        </label>
        <label v-if="page_number+1<=number_of_pages" class="btn btn-dark" @click="getStream(page_number+1)">
            next page
        </label>
    </div>
</template>