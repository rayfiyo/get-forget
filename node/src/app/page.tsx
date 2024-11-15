// src/app/page.tsx

import React from 'react';
import GetForget from '@/components/get-forget/GetForget';

const Page = () => {
  return (
    <main className="flex flex-col items-center justify-center min-h-screen bg-gray-50 p-4">
      <h1 className="text-2xl font-bold mb-4">Get Forget</h1>
      <GetForget />
    </main>
  );
};

export default Page;

