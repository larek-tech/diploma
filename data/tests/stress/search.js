import http from 'k6/http';
import { check } from 'k6';

const { HOSTNAME, PORT, SOURCE_ID } = require('./config.js');


const url = `http://${HOSTNAME}:${PORT}/q`;

export const options = {
    stages: [
        { duration: '5m', target: 10 },
        { duration: '5m', target: 20 },
        { duration: '5m', target: 30 },
        { duration: '5m', target: 40 },
        { duration: '5m', target: 50 },
        { duration: '5m', target: 60 },
        { duration: '5m', target: 0 },// ramp-down to 0 users
    ],
    httpDebug: 'none',
    logFormat: 'json',
    noVUConnectionReuse: false,
    discardResponseBodies: true,
};
const queries = [
    {
        query: "что такое gitlic",
        sourceIds: [SOURCE_ID],
        topK: 3,
        threshold: 0.1,
        useQuestions: false
    },
    {
        query: "пример второго запроса",
        sourceIds: [SOURCE_ID],
        topK: 5,
        threshold: 0.2,
        useQuestions: true
    },
    // Add more queries as needed
];

const params = {
    headers: {
        'Content-Type': 'application/json',
    },
};

export default function () {
    // Pick a query based on the virtual user or iteration
    const idx = __ITER % queries.length;
    const payload = JSON.stringify(queries[idx]);
    let res = http.post(url, payload, params);
    check(res, { 'status was 200': (r) => r.status == 200 });
}