/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_MINIO_BASE?: string;
  /** База для относительных имён файлов фото/видео, как в шаблоне: http://localhost:9000/test */
  readonly VITE_MEDIA_BASE?: string;
}

declare module "*.mp4" {
  const src: string;
  export default src;
}

declare module "*.svg" {
  const src: string;
  export default src;
}

declare module "*.jpg" {
  const src: string;
  export default src;
}

declare module "*.png" {
  const src: string;
  export default src;
}
