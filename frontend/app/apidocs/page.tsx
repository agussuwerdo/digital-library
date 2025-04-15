'use client';

import { RedocStandalone } from 'redoc';
import { useEffect, useState } from 'react';

interface PathItem {
  get?: Operation;
  post?: Operation;
  put?: Operation;
  delete?: Operation;
  parameters?: Parameter[];
}

interface Operation {
  tags?: string[];
  summary?: string;
  description?: string;
  parameters?: Parameter[];
  requestBody?: {
    content: {
      [key: string]: {
        schema: Schema;
      };
    };
  };
  responses: {
    [key: string]: {
      description: string;
      content?: {
        [key: string]: {
          schema: Schema;
        };
      };
    };
  };
}

interface Parameter {
  name: string;
  in: string;
  description?: string;
  required?: boolean;
  schema: Schema;
}

interface Schema {
  type: string;
  format?: string;
  items?: Schema;
  properties?: {
    [key: string]: Schema;
  };
  required?: string[];
  $ref?: string;
}

interface OpenAPISpec {
  openapi: string;
  info: {
    title: string;
    version: string;
    description?: string;
  };
  paths: {
    [key: string]: PathItem;
  };
  components?: {
    schemas?: {
      [key: string]: Schema;
    };
    securitySchemes?: {
      [key: string]: {
        type: string;
        scheme?: string;
        bearerFormat?: string;
      };
    };
  };
}

export default function ApiDocPage() {
  const [spec, setSpec] = useState<OpenAPISpec | null>(null);

  useEffect(() => {
    const fetchSpec = async () => {
      try {
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/apidocs`);
        const data = await response.json();
        setSpec(data);
      } catch (error) {
        console.error('Error fetching API spec:', error);
      }
    };

    fetchSpec();
  }, []);

  if (!spec) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-t-2 border-b-2 border-gray-900"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <RedocStandalone
        spec={spec}
        options={{
          nativeScrollbars: true,
          theme: {
            colors: {
              primary: {
                main: '#3B82F6'
              }
            },
            typography: {
              fontFamily: 'Inter, system-ui, sans-serif',
              headings: {
                fontFamily: 'Inter, system-ui, sans-serif',
              }
            },
            sidebar: {
              width: '300px',
            }
          }
        }}
      />
    </div>
  );
} 