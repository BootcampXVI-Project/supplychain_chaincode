import {Object, Property} from 'fabric-contract-api';

@Object()
export class CounterNO {

    @Property()
    public docType?: string;

    @Property()
    public Counter: number;

}
