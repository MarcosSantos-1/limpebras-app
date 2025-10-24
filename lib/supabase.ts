// Cliente Supabase temporário para substituir dependências
export const supabase = {
  auth: {
    getUser: async () => {
      console.log('Auth getUser simulado');
      return {
        data: { 
          user: { 
            id: 'mock-user-id', 
            email: 'mock@example.com' 
          } 
        },
        error: null
      };
    }
  },
  storage: {
    from: (bucket: string) => ({
      upload: async (path: string, file: File) => {
        console.log('Upload simulado:', { bucket, path, file });
        return {
          data: { path },
          error: null as any // Permitir acesso a propriedades do erro
        };
      },
      remove: async (paths: string[]) => {
        console.log('Remove simulado:', { bucket, paths });
        return { error: null as any };
      },
      getPublicUrl: (path: string) => {
        console.log('GetPublicUrl simulado:', { bucket, path });
        return {
          data: { publicUrl: `https://mock-supabase.com/storage/v1/object/public/${bucket}/${path}` }
        };
      }
    })
  },
  from: (table: string) => ({
    insert: async (data: any[]) => {
      console.log('Insert simulado:', { table, data });
      return { error: null as any };
    },
    select: () => ({
      eq: () => ({
        single: async () => {
          console.log('Select simulado:', { table });
          return { data: null, error: null as any };
        }
      })
    }),
    update: () => ({
      eq: () => ({
        single: async () => {
          console.log('Update simulado:', { table });
          return { data: null, error: null as any };
        }
      })
    }),
    delete: () => ({
      eq: () => ({
        single: async () => {
          console.log('Delete simulado:', { table });
          return { data: null, error: null as any };
        }
      })
    })
  })
};

