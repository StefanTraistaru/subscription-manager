export class Subscription {
    name: string;
    details: string;
    price: string;
    date_d: string;
    date_m: string;
    date_y: string;

    constructor(name: string, details: string, price: string, date_d: string, date_m: string, date_y: string) {
        this.name = name;
        this.details = details;
        this.price = price;
        this.date_d = date_d;
        this.date_m = date_m;
        this.date_y = date_y;
    }
}