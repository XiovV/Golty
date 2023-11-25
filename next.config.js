/** @type {import('next').NextConfig} */
const nextConfig = {
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
