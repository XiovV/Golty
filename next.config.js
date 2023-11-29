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
    ];
  },
};

module.exports = nextConfig;
