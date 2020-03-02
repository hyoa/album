describe('endpoint /album/{slug}', () => {
  beforeEach(() => {
    cy.login()

    cy.visit('/#/album/album-5')
    cy.wait(2000)
  })

  it('Visit the endpoint /album/{slug} should display the album', () => {
    cy.get('[data-e2e="album-title"]').should('have.length', 1).and('be.visible')
    cy.get('[data-e2e="album-description"]').should('have.length', 1)
    cy.get('[data-e2e="album-informations"]').should('have.length', 1).and('be.visible')
    cy.get('[data-e2e="medias-grid"]').should('have.length', 1).and('be.visible')
  })

  it('Clicking on a media should open the gallery', () => {
    cy.get('[data-e2e="medias-grid"] img').click()
    cy.wait(200)

    cy.get('.slide').should('have.length', 1)
    cy.get('.slide-content').first().should('be.visible')
  })
})
