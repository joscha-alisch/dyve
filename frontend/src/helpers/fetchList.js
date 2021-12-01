import axios from "axios";


export const fetchList = (apiName) => (perPage, page, setResults) => {
    axios.get("/api/" + apiName + "?perPage=" + perPage + "&page=" + page)
        .then((res) => {
            if (res.data.result) {
                setResults(res.data.result[apiName], res.data.result.totalResults)
            }
        })
}