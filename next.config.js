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
  async redirects() {
    return [
      {
        source: "/",
        destination: "/channels",
        permanent: true,
      },
    ];
  },
};

module.exports = nextConfig;
