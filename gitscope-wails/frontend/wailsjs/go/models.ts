export namespace git {
	
	export class RemoteInfo {
	    Name: string;
	    FetchURL: string;
	    PushURL: string;
	
	    static createFrom(source: any = {}) {
	        return new RemoteInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.FetchURL = source["FetchURL"];
	        this.PushURL = source["PushURL"];
	    }
	}

}

