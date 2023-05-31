import { FileUploadWithPreview } from 'file-upload-with-preview'
import 'file-upload-with-preview/dist/style.css';

let upload: FileUploadWithPreview | null = null

export const randomString = (length: number) => (Math.random() + 1).toString(36).substring(length);

export const createFileUpload = (hash: string) => {
    upload = new FileUploadWithPreview(hash);

    console.log('FUCK')

    return hash
}