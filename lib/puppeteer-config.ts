// Configuração do Puppeteer para geração de PDFs
export function getPuppeteerConfig() {
  return {
    headless: true,
    args: [
      '--no-sandbox',
      '--disable-setuid-sandbox',
      '--disable-dev-shm-usage',
      '--disable-accelerated-2d-canvas',
      '--no-first-run',
      '--no-zygote',
      '--disable-gpu',
      '--disable-web-security',
      '--disable-features=VizDisplayCompositor'
    ]
  };
}

export default getPuppeteerConfig;
