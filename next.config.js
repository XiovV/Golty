/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "i.imgur.com",
        port: "",
        // pathname: "/",
      },
      {
        protocol: "https",
        hostname: "yt3.googleusercontent.com",
        port: "",
        // pathname: "/",
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
