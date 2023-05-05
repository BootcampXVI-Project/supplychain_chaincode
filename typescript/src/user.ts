import {Object, Property} from 'fabric-contract-api';

@Object()
export class User {
    @Property()
    public docType?: string;

    @Property()
    public UserId: string;

    @Property()
    public Email: string;

    @Property()
    public Password: string;

    @Property()
    public UserName: string;

    @Property()
    public Address: string;
    
    @Property()
    public UserType: string;

    @Property()
    public Role?: string;

    @Property()
    public Status: string;
}
