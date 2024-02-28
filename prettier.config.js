/** @type {import('prettier').Config} */
module.exports = {
  trailingComma: 'es5',
  tabWidth: 2,
  semi: false,
  singleQuote: true,

  overrides: [
    {
      files: '.postcssrc',
      options: { parser: 'json' },
    },
  ],
}
