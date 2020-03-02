// https://docs.cypress.io/api/introduction/api.html

describe('endpoint / without auth', () => {
  it('Visit the endpoint / without login should redirect to auth page', () => {
    cy.visit('/')
    cy.url().should('contain', '/auth')
  })
})

describe('endpoint / with auth', () => {
  beforeEach(() => {
    cy.login()
    cy.visit('/')
  })

  it('Visit the endpoint / logged should bring to homepage', () => {
    cy.url().should('not.contain', '/auth')

    cy.wait(500)

    cy.get('[data-e2e="search-form-section"]').should('have.length', 1)
    cy.get('[data-e2e="last-albums-section"]').should('have.length', 1)

    cy.get('[data-e2e="album-card"]').should('have.length', 3)
    cy.get('[data-e2e="album-card"] img').should('exist')
  })

  it('Search should display matching albums', () => {
    cy.get('[data-e2e="search-form-section"] input').type('album')
    cy.get('[data-e2e="search-form-section"] button').click()

    cy.wait(2000)

    cy.get('[data-e2e="search-result-section"]').should('have.length', 1)
    cy.get('[data-e2e="album-card"]').should('exist')
  })

  it('Clicking on a card should redirect to the album', () => {
    cy.get('[data-e2e="album-card"]').first().click()

    cy.url().should('contain', '/album/')
  })
})
