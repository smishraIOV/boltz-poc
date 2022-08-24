
export default class BoltzPocClient {

    static async CreatePayment() : Promise<string> {
        const request = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ amount: 23450 })
        };
        return fetch('http://localhost:8080/payment', request)
            .then(response => response.json())
            .then(data => data.lninvoice);
    }
}