import {Object, Property} from 'fabric-contract-api';
import {Product} from './product';

@Object()
export class ProductHistory {
    @Property()
    public docType?: string;

    @Property()
    public Record: Product;

    @Property()
    public TxId: string;

    @Property()
    public Timestamp: string; // time

    @Property()
    public IsDelete: boolean;

}
