import { Subscription } from "./subscription";

export class ResponseResult {
    error: string;
    result: string;
    data: Subscription[];

    constructor(error: string, result: string, data: Subscription[]) {
        this.error = error;
        this.result = result;
        this.data = data;
    }
}