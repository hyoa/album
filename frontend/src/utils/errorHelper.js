export default code => {
  switch (code) {
    case 11:
      return 'Cet adresse email est déjà utilisé !'
    case 20:
      return 'Cet album existe déjà'
    case 21:
      return 'Le titre ne peut pas être vide !'
    default:
      return 'Quelque chose ne s\'est pas déroulé comme il faut !'
  }
}
