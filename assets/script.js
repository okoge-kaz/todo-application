/* placeholder file for JavaScript */
const confirm_delete = ((id) => {
  if (window.confirm('Are you sure you want to delete this?')) {
    location.href = `/task/${id}/delete`;
  }
})

const confirm_update = ((id) => {
  if (window.confirm('Are you sure you want to update this?')) {
    location.href = `/task/${id}/edit`;
  }
})
