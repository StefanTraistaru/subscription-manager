export class User {
    username: string;
    firstname: string;
    lastname: string;
    password: string;
    token: string;

    constructor(username: string, firstname: string, lastname: string, password: string, token: string) {
        this.username = username;
        this.firstname = firstname;
        this.lastname = lastname;
        this.password = password;
        this.token = token;
    }
}