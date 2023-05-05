import {Object, Property} from 'fabric-contract-api';

@Object()
export class ProductActor {
    
    @Property()
    public docType?: string;

    @Property()
    public SupplierId: string;

    @Property()
    public ManufacturerId: string;

    @Property()
    public DistributorId: string;

    @Property()
    public Retailer: string;

}
