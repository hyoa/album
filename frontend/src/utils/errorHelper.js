export default code => {
  switch (code) {
    case 11:
      return 'error.emailAlreadyUsed'
    case 20:
      return 'error.albumAlreadyExist'
    case 21:
      return 'error.titleCannotBeEmpty'
    default:
      return 'error.default'
  }
}
