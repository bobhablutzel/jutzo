import {fetchData} from "@/services/fetchData";

function retrieveBlogEntry(id, blogEntrySetter) {
    const blogURL = `http://192.168.1.16:8081/blog/entry/${id}`
    fetchData(blogURL,(data, error) => {
            blogEntrySetter( error ? buildErrorEntry(error) : data)
        }
    )
}

function retrieveNewest( setter ) {
    const newestURL = `http://192.168.1.16:8081/blog/newest`
    fetchData(newestURL, (data, error ) => {
        setter( error ? buildErrorEntry(error) : data )
    })
}

function buildErrorEntry(error) {
    return {
        title: 'An error occurred',
        body: [ { text: error.message } ],
        teaser: error.message,
        id: '00000000-0000-0000-0000-00000000000',
        publicationDate: Date.now()
    }
}



export { retrieveBlogEntry, retrieveNewest }