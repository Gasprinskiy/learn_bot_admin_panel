import antfu from '@antfu/eslint-config';

export default antfu({
  extends: ['plugin:vuejs-accessibility/recommended'],
  stylistic: {
    semi: true,
    quotes: 'single',
    arrowParens: true,
    braceStyle: '1tbs',
    allowSingleLine: false,
    blockSpacing: true,
    commaDangle: 'never',
    overrides: {
      'curly': 'error',
      'style/brace-style': 'error',
      'ts/no-dynamic-delete': 'error',
      'perfectionist/sort-imports': 'off',
      'perfectionist/sort-named-imports': 'off',
      'perfectionist/sort-exports': 'off',
      'unicorn/consistent-function-scoping': 'off',
      'unicorn/new-for-builtins': 'off',
    },
  },
  yaml: false,
  rules: {
    '@typescript-eslint/consistent-type-definitions': 'off',
  },
});
