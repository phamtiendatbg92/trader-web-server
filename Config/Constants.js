
const Constants = {
    SERVER_URL: "localhost:5000",
    PUBLIC_IMG: "/public/images",
    SERVER_PUBLIC_KEY : "SERVER_PUBLIC_KEY", // this variable is used to replace later in client
}

// map url to ID in database
const UrlMap = new Map();

module.exports = { Constants, UrlMap }


