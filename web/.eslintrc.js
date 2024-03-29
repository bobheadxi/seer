module.exports = {
  root: true,
  env: {
    node: true
  },
  extends: [
    'plugin:vue/essential',
    '@vue/airbnb',
    '@vue/typescript',
  ],
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-extra-semi': 'off',
    'max-len': 'off',
    'lines-between-class-members': 'off',
    'class-methods-use-this': 'off',
    'import/prefer-default-export': 'off',
  },
  parserOptions: {
    parser: '@typescript-eslint/parser'
  }
}
