export declare const randomString: (length: number) => string;
export declare const createFileUpload: (hash: string) => string;
export declare const createFileChunks: (alreadyUploaded: string) => Promise<Int8Array[]>;
