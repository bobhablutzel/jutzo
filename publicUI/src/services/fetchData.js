// fetchData is the main function for getting data from the server
// This uses Go-style returns; the first element is the JSON
// payload returned while the second is an error (if any, or
// null for success)
export function fetchData(url, callback) {
    fetch(url, {
        method: 'GET',
        headers: {
            'content-type': 'application/json'
        }
    })
        .then(res => {
            // a non-200 response code
            if (!res.ok) {
                // create error instance with HTTP status text
                const error = new Error(res.statusText);
                error.json = res.json();
                callback(null, error)
            } else {
                res.json().then(json => {
                    callback(json, null)
                })
            }
        })
        .catch((error) => {
            callback(null, error)
        })
}