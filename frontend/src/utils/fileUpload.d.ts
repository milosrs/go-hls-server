export declare const randomString: (length: number) => string;
export declare const createFileUpload: (hash: string) => string;
export declare const createFileChunks: (alreadyUploaded: string) => Promise<Uint8Array[]>;
export declare const sendFileChunks: (chunks: Uint8Array[]) => Promise<void>;
