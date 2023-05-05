import {Object, Property} from 'fabric-contract-api';

@Object()
export class ProductDate {
    
    @Property()
    public docType?: string;

    @Property()
    public Cultivated: string;

    @Property()
    public Harvested: string;

    @Property()
    public Imported: string;

    @Property()
    public Manufacturered: string;

    @Property()
    public Exported: string;
    
    @Property()
    public Distributed: string;

    @Property()
    public Sold?: string;

}
