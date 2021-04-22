const express = require('express');
const app = express();

app.get('/', async (req, res) => {
    res.send({ message: "hello!" });
});

app.listen(3001, () => {
    console.log('listening on port 3001');
});