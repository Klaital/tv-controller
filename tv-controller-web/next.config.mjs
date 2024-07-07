const isProd = process.env.NODE_ENV === 'production'

/** @type {import('next').NextConfig} */
const nextConfig = {
    output: "export",
    distDir: "dist",
    basePath: isProd ? "/web" : undefined,
};

export default nextConfig;
