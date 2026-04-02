export namespace main {
	
	export class SendMessageResponse {
	    content: string;
	    commandID: string;
	
	    static createFrom(source: any = {}) {
	        return new SendMessageResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = source["content"];
	        this.commandID = source["commandID"];
	    }
	}

}

export namespace session {
	
	export class SuggestedCommand {
	    id: string;
	    session_id: string;
	    terminal_id: string;
	    command: string;
	    description: string;
	    context: string;
	    // Go type: time
	    created_at: any;
	    used_count: number;
	    // Go type: time
	    last_used_at: any;
	
	    static createFrom(source: any = {}) {
	        return new SuggestedCommand(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.session_id = source["session_id"];
	        this.terminal_id = source["terminal_id"];
	        this.command = source["command"];
	        this.description = source["description"];
	        this.context = source["context"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.used_count = source["used_count"];
	        this.last_used_at = this.convertValues(source["last_used_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

