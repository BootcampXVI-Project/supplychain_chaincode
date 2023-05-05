import {Object, Property} from 'fabric-contract-api';
import {ProductDate} from './productDate';
import {ProductActor} from './productActor';

@Object()
export class Product {
    
    @Property()
    public docType?: string;

    @Property()
    public ProductId: string;

    @Property()
    public ProductName: string;

    @Property()
    public Dates: ProductDate;

    @Property()
    public Actors: ProductActor;

    @Property()
    public Price: number;
    
    @Property()
    public Status: string;

    @Property()
    public Description: string;

}
