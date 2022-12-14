import {fetchJSON} from "@/services/fetch.js";

function retrieveBlogEntry(id, blogEntrySetter) {
    const blogURL = `${process.env.VUE_APP_API_URL}v1/blog/entry/${id}`
    fetchJSON(blogURL,(data, error) => {
            blogEntrySetter( error ? buildErrorEntry(error) : data)
        }
    )
}

function retrieveNewest( setter ) {
    const newestURL = `${process.env.VUE_APP_API_URL}v1/blog/newest`
    fetchJSON(newestURL, (data, error ) => {
        setter( error ? buildErrorEntry(error) : data )
    })
}

function buildErrorEntry(error) {
    return [ {
        title: 'An error occurred',
        body: [ { text: error.message } ],
        teaser: error.message,
        id: '00000000-0000-0000-0000-00000000000',
        publicationDate: new Date().toISOString()
    } ]
}



export { retrieveBlogEntry, retrieveNewest }