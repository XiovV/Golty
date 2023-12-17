/** @type {import('next').NextConfig} */
const nextConfig = {
  eslint: {
    ignoreDuringBuilds: true,
  },
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "i.imgur.com",
        port: "",
      },
      {
        protocol: "https",
        hostname: "yt3.googleusercontent.com",
        port: "",
      },
      {
        protocol: "https",
        hostname: "i.ytimg.com",
        port: "",
      },
    ],
  },
  pageExtensions: ["ts", "tsx"],
  async redirects() {
    return [
      {
        source: "/dashboard",
        destination: "/dashboard/channels",
        permanent: true,
      },
      {
        source: "/",
        destination: "/dashboard/channels",
        permanent: true,
      },
    ];
  },
};

module.exports = nextConfig;
